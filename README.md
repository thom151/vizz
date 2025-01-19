# Vizz - Visualizing Your Stories

##  📖 Description
**Vizz** is an interactive website that would transform stories into images. By uploading a story or a book, users will be generated images per chapter as they read.

## 🤔 Why Vizz?
I love reading books. Reading books requires a lot of imagination. However, I often find it challenging to visualize intricate plots or alien worlds, especally reading science fiction books. Viz was created to bridge the gap between imaging and storytelling, offering readers a way to see the stories come to life.


---
## 🚀 Quick Start
Follow these steps to set up and run Viz on your local machine.
Before starting, ensure you have:
1. [Go](https://go.dev/) installed on your machine.
2. A `.env` file in the root directory with the following environment variables configured:
   ```env
   DATABASE_URL=libsql://vizz-db-...  # Your Turso database connection
   PLATFORM=dev                       # Set to 'dev' for local development
   SECRET=XcjOa...                    # Secret key for tokens
   PORT=8080                          # Port the server will run on
   OPEN_API=sk-Hgtk...                # OpenAI API key
   ASSISTANT=asst_gA4...              # Assistant functionality (optional)
   ```

### How to Obtain Environment Variables
1. **DATABASE_URL**
   - Create a Turso database account at [https://turso.tech](https://turso.tech) and set up a new database.
   - Click create group and then click aws.
   - Click continue with AWS.
   - Enter a name, select a region and then click `Create Group`.
   - Click the database url to copy. It should look like `dbName-username.aws-us-east-1.turso.io`.
   - On the right side, click the three dots again and click `Create Token`.
   - Lastly, copy your database url and connect it with `?authToken={your newly created token}`.
  
2. **SECRET**
   - Generate a secret key by running the following command in your terminal:
     ```bash
     openssl rand -base64 64
     ```
   - Copy the output and paste it into your `.env` file as the `SECRET` value.

3. **OPEN_API**
   - Sign up for an OpenAI account at [https://platform.openai.com/signup](https://platform.openai.com/signup).
   - Go to the API Keys section in your OpenAI dashboard and generate a new API key.
   - Copy the key and paste it into your `.env` file as the `OPEN_API` value.

4. **ASSISTANT**
   - Past the following: asst_gA4LrW3d74SvGeZjW8gHnPfv
---
### Steps to Run Locally

1. **Clone the Repository**
    Clone the repository to your local machine:
    ```bash
    git clone https://github.com/thom-151/vizz.git
    cd vizz

2. **Install Dependencies**
    Download the required Go modules:
    go mod tidy

3. **Set up Environment Variables**
    Create a `.env` file in the project root directory and configure the variables as shown in the example above.               

4. **Migrate the Database**  
   Run the migration script to set up your database:
   ```bash
   ./scripts/migrateup.sh
   ```
5. **Run the Application**  
   Start the server with the following command:
   ```bash
   go run main.go
   ```
6. **Access the Application**  
   Open your browser and navigate to:
   ```
   http://localhost:8080

    ---



