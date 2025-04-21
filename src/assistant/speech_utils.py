import os
import tempfile
import logging
from google.cloud import speech
import ffmpeg

logger = logging.getLogger(__name__)

async def transcribe_ogg_voice(file_path: str, primary_language: str = None) -> str:
    """Convert OGG to FLAC and transcribe using Google Cloud Speech-to-Text.
    
    Supports automatic language detection between Ukrainian and English.
    
    Args:
        file_path: Path to the OGG voice file
        primary_language: Optional primary language code ("uk-UA" or "en-US")
                          If not provided, will use auto language detection
    
    Returns:
        Transcribed text
    """
    # Convert OGG (Opus) to FLAC using ffmpeg
    with tempfile.NamedTemporaryFile(suffix=".flac", delete=False) as flac_file:
        flac_path = flac_file.name
    try:
        (
            ffmpeg
            .input(file_path)
            .output(flac_path, ac=1, ar='16000', format='flac')
            .run(quiet=True, overwrite_output=True)
        )
        
        # Transcribe
        client = speech.SpeechClient()
        with open(flac_path, "rb") as audio_file:
            audio_content = audio_file.read()
        
        audio = speech.RecognitionAudio(content=audio_content)
        
        # If primary language specified, we'll use it with alternatives
        # Otherwise, use automatic language detection
        if primary_language:
            # Try with specified primary language and alternative languages
            config = speech.RecognitionConfig(
                encoding=speech.RecognitionConfig.AudioEncoding.FLAC,
                sample_rate_hertz=16000,
                language_code=primary_language,  # Primary language
                alternative_language_codes=["uk-UA" if primary_language == "en-US" else "en-US"],
                model="default"
            )
        else:
            # Try with automatic language detection
            config = speech.RecognitionConfig(
                encoding=speech.RecognitionConfig.AudioEncoding.FLAC,
                sample_rate_hertz=16000,
                language_code="uk-UA",  # Primary language 
                alternative_language_codes=["en-US"],  # Alternative language
                model="default"
            )
        
        logger.info(f"Transcribing with config: {config}")
        response = client.recognize(config=config, audio=audio)
        
        if not response.results:
            logger.warning("No transcription results returned")
            return ""
            
        transcript = " ".join([result.alternatives[0].transcript for result in response.results])
        detected_language = response.results[0].language_code if response.results else "unknown"
        logger.info(f"Detected language: {detected_language}")
        
        return transcript.strip()
    except Exception as e:
        logger.error(f"Error in transcription: {e}")
        return f"[Transcription error: {e}]"
    finally:
        # Clean up the temporary file
        if os.path.exists(flac_path):
            os.remove(flac_path)
