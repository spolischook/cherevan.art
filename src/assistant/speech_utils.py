import os
import tempfile
from google.cloud import speech
import ffmpeg

async def transcribe_ogg_voice(file_path: str) -> str:
    """Convert OGG to FLAC and transcribe using Google Cloud Speech-to-Text."""
    # Convert OGG (Opus) to FLAC using ffmpeg
    with tempfile.NamedTemporaryFile(suffix=".flac", delete=False) as flac_file:
        flac_path = flac_file.name
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
    config = speech.RecognitionConfig(
        encoding=speech.RecognitionConfig.AudioEncoding.FLAC,
        sample_rate_hertz=16000,
        language_code="en-US"
    )
    response = client.recognize(config=config, audio=audio)
    transcript = " ".join([result.alternatives[0].transcript for result in response.results])
    os.remove(flac_path)
    return transcript.strip()
