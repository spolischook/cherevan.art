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
    # Download the voice file from Telegram
    file = await context.bot.get_file(update.message.voice.file_id)
    with tempfile.NamedTemporaryFile(suffix=".ogg", delete=False) as ogg_file:
        ogg_path = ogg_file.name
        async with aiohttp.ClientSession() as session:
            async with session.get(file.file_path) as resp:
                ogg_file.write(await resp.read())
    # Transcribe audio
    transcript = await transcribe_ogg_voice(ogg_path)
    # Compose prompt for Gemini
    prompt = transcript
    # Get Gemini response
    try:
        gemini_response = gemini_generate_response(prompt)
    except Exception as e:
        logger.error(f"Gemini API error: {e}")
        gemini_response = "[Error: Could not get response from Gemini]"
    # Send three messages
    user_display = user.full_name if user else str(chat_id)
    await update.message.reply_text(f"User: {user_display}")
    await update.message.reply_text(transcript or "[No transcription]")
    await update.message.reply_text(gemini_response)

async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Sends a welcome message when the /start command is issued."""
    await update.message.reply_text('Hello! Send me a message, and I will echo it back.')

async def echo(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Echoes the user message."""
    if update.message and update.message.text:
        logger.info(f"Received message from {update.effective_chat.id}: {update.message.text}")
        # Simple dummy response: send the same text back
        await update.message.reply_text(update.message.text)
    else:
        logger.info(f"Received non-text message from {update.effective_chat.id}")
        await update.message.reply_text("I can only echo text messages right now.")


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

    # Register message handler for text messages (but not commands)
    application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, echo))

    # Register handler for voice messages
    application.add_handler(MessageHandler(filters.VOICE, handle_voice))

    # Start the Bot using polling - checks Telegram for new messages periodically.
    # This is good for development but not ideal for serverless deployment.
    logger.info("Starting bot using polling...")
    application.run_polling()

if __name__ == '__main__':
    main()
