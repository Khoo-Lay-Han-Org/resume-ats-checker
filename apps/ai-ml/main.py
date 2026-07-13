import uvicorn
from api.config import app
import api.view

if __name__ == "__main__":
    uvicorn.run(app, host="localhost", port=9000)
