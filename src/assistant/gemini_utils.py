import os
from google.cloud import aiplatform

def gemini_generate_response(prompt: str) -> str:
    """
    Sends a prompt to Gemini via Vertex AI and returns the response text.
    
    Uses the Google Cloud service account credentials (same as speech recognition).
    No additional API key is required as it uses the credentials in your JSON file.
    
    Supported Gemini models on Vertex AI (as of 2024):
    - gemini-1.0-pro
    - gemini-1.5-flash
    - gemini-1.5-pro
    """
    # You can set these as environment variables or hardcode them here
    model_name = os.getenv("GEMINI_MODEL", "gemini-1.5-flash")
    project_id = os.getenv("GCP_PROJECT", "") 
    location = os.getenv("GCP_LOCATION", "us-central1")
    
    try:
        # Initialize Vertex AI with default credentials (uses service account JSON)
        aiplatform.init(project=project_id, location=location)
        
        # Import VertexAI prediction service for Gemini
        # This needs to be imported here to avoid circular import issues
        from vertexai.generative_models import GenerativeModel
        
        # Create a model instance
        model = GenerativeModel(model_name)
        
        # Generate content
        response = model.generate_content(prompt)
        
        # Return the text response
        return response.text
    except Exception as e:
        import traceback
        print(f"Vertex AI Gemini error: {e}")
        traceback.print_exc()
        return f"[Error: Could not get response from Gemini: {e}]"
