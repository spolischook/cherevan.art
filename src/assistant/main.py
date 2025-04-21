import logging
import os
from dotenv import load_dotenv

from telegram import Update
from telegram.ext import Application, CommandHandler, MessageHandler, filters, ContextTypes
import tempfile
import aiohttp
from gemini_utils import gemini_generate_response
from speech_utils import transcribe_ogg_voice

# --- Configuration ---
# Load environment variables from a .env file (for local development)
load_dotenv()
TELEGRAM_BOT_TOKEN = os.getenv("TELEGRAM_BOT_TOKEN")

# Enable logging
logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    level=logging.INFO
)
logger = logging.getLogger(__name__)

# --- Bot Handlers ---

async def handle_voice(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Handles incoming voice messages: downloads, transcribes, sends to Gemini, and replies."""
    user = update.effective_user
    chat_id = update.effective_chat.id
    if not update.message or not update.message.voice:
        await update.message.reply_text("No voice message found.")
        return
    
    logger.info(f"Received voice message from {chat_id}")
    
    # Get user language preference if available
    user_data = context.user_data
    primary_language = user_data.get('primary_language') if user_data else None
    
    # Download the voice file from Telegram
    file = await context.bot.get_file(update.message.voice.file_id)
    with tempfile.NamedTemporaryFile(suffix=".ogg", delete=False) as ogg_file:
        ogg_path = ogg_file.name
        async with aiohttp.ClientSession() as session:
            async with session.get(file.file_path) as resp:
                ogg_file.write(await resp.read())
    
    # Status message to show processing
    status_msg = await update.message.reply_text("🎙️ Transcribing your message...")
    
    # Transcribe audio with automatic language detection (Ukrainian/English)
    transcript = await transcribe_ogg_voice(ogg_path, primary_language)
    
    # Clean up the temporary file
    if os.path.exists(ogg_path):
        os.remove(ogg_path)
    
    # Compose prompt for Gemini
    prompt = transcript
    
    # Update status message
    await status_msg.edit_text("🤖 Getting AI response...")
    
    # Get Gemini response
    try:
        gemini_response = gemini_generate_response(prompt)
    except Exception as e:
        logger.error(f"Gemini API error: {e}")
        gemini_response = "[Error: Could not get response from Gemini]"
    
    # Delete the status message
    await status_msg.delete()
    
    # Send messages with results
    user_display = user.full_name if user else str(chat_id)
    await update.message.reply_text(f"User: {user_display}")
    await update.message.reply_text(transcript or "[No transcription]")
    await update.message.reply_text(gemini_response)

async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Sends a welcome message when the /start command is issued."""
    await update.message.reply_text('Hello! I can transcribe voice messages in English and Ukrainian. Use /language to set your preferred language.')

async def language_command(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Handles the /language command to set user language preference."""
    args = context.args
    
    if not args or args[0].lower() not in ['en', 'uk', 'auto']:
        await update.message.reply_text(
            "Please specify a language code: /language [code]\n\n"
            "Available options:\n"
            "- /language en - Prioritize English\n"
            "- /language uk - Prioritize Ukrainian\n"
            "- /language auto - Auto detect (default)"
        )
        return
    
    lang_code = args[0].lower()
    
    # Map language code to full code with regional variant
    language_map = {
        'en': 'en-US',
        'uk': 'uk-UA',
        'auto': None  # None means auto-detection
    }
    
    # Store in user data
    context.user_data['primary_language'] = language_map[lang_code]
    
    # Confirm setting
    if lang_code == 'auto':
        await update.message.reply_text("Language preference set to: Auto detect (Ukrainian + English)")
    else:
        await update.message.reply_text(f"Language preference set to: {lang_code.upper()}")

async def handle_text(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Handles text messages by sending them to Gemini and returning the response."""
    if update.message and update.message.text:
        user = update.effective_user
        chat_id = update.effective_chat.id
        logger.info(f"Received text message from {chat_id}: {update.message.text}")
        
        # Get the message text
        prompt = update.message.text
        
        # Show typing indicator
        await update.message.chat.send_action(action="typing")
        
        # Get Gemini response
        try:
            gemini_response = gemini_generate_response(prompt)
        except Exception as e:
            logger.error(f"Gemini API error: {e}")
            gemini_response = "[Error: Could not get response from Gemini]"
        
        # Send response
        await update.message.reply_text(gemini_response)
    else:
        logger.info(f"Received empty text message from {update.effective_chat.id}")
        await update.message.reply_text("Please send a non-empty text message.")


# --- Main Execution ---

def main() -> None:
    """Start the bot."""
    if not TELEGRAM_BOT_TOKEN:
        logger.error("TELEGRAM_BOT_TOKEN environment variable not set!")
        return

    # Create the Application and pass it your bot's token.
    application = Application.builder().token(TELEGRAM_BOT_TOKEN).build()

    # Register command handlers
    application.add_handler(CommandHandler("start", start))
    application.add_handler(CommandHandler("language", language_command))

    # Register message handler for text messages (but not commands)
    application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, handle_text))

    # Register handler for voice messages
    application.add_handler(MessageHandler(filters.VOICE, handle_voice))

    # Start the Bot using polling - checks Telegram for new messages periodically.
    # This is good for development but not ideal for serverless deployment.
    logger.info("Starting bot using polling...")
    application.run_polling()

if __name__ == '__main__':
    main()
