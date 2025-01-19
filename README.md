# Vizz - Visualizing Your Stories

##  📖 Description
**Vizz** is an interactive website that would transform stories into images. By uploading a story or a book, users will be generated images per chapter as they read.

## 🤔 Why Vizz?
I have recently been reading a lot of science-fiction books. Reading books requires a lot of imagination. However, I often find it challenging to visualize intricate plots or alien worlds, especally reading science fiction books. Viz was created to bridge the gap between imaging and storytelling, offering readers a way to see the stories come to life.


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
   - **Note**: Using the OpenAI API for image generation will incur costs. Allocate a budget of around $1–$2 to cover the API usage fees, depending on the number of images generated.

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

4. **Install Goose for Database Migrations**  
   Install the `goose` CLI tool for managing database migrations:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```
5. **Migrate the Database**  
   Run the migration script to set up your database:
   ```bash
   ./scripts/migrateup.sh
   ```
6. **Run the Application**  
   Start the server with the following command:
   ```bash
   go run .
   ```
7. **Access the Application**  
   Open your browser and navigate to:
   ```
   http://localhost:8080

    ---
## 🛠 Usage

Here’s how to interact with the website once it’s running locally:
### **Creating an account**
1. Change the url from `http://localhost:8080` to `http://localhost:8080/api/users`.
2. Create an account using any dummy email or password.

### **Login**
1. Open your browser and navigate to `http://localhost:8080`.
2. If you are logged in, you'll be able to click the buttons.
2. If not, you need to log in.
3. Enter your email and password, then submit the form to log in.

### **Upload a Book**
1. Once logged in, navigate to the **Upload Book** section.
2. Use the provided form to select an ePub file from your computer and upload it.
3. The book will be processed and stored in your account.

### **Search for a Story**
1. Go to the **Search** section on the website.
2. Enter the title of the book you want to find in the search bar.
3. View the results and select the desired book to proceed.

### **Visualize a Story**
1. Click on a book from the search results or your uploaded books list.
2. The website will display the visualized story, with options for pagination.
3. Use the navigation buttons to move between pages or sections of the story.

---

## 📚 Features

- **Login**: Authenticate with email and password to access the application.
- **Upload Books**: Easily upload text files to generate visualizations.
- **Search Functionality**: Look up books or stories by title.
- **Story Visualization**: Transform uploaded stories into rich visual representations, with pagination and OpenAI-powered threads.

---

## 🔒 Notes on Security

1. **Environment Separation**: Use `PLATFORM=dev` for development and `PLATFORM=prod` for stricter production setups.
2. **Sensitive Data**: Protect sensitive environment variables by keeping the `.env` file secure and out of version control.
3. **OpenAI API Costs**: Be mindful of the API usage as it may incur costs based on usage.

---
## 🤝 Contributing

### Clone the repo

```bash
git clone https://github.com/thom151/vizz
cd vizz
```

### Build the project

```bash
go build
```

### Run the project
```bash
./vizz
```

### Run the tests

```bash
go test ./...
```

---

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

## 🤝 Need Help?

If you encounter any issues or have questions:
- Open an issue in this repository.
- Contact me at [thomassantos2003@gmail.com](mailto:thomassantos2003@gmail.com)

---

Enjoy using **Vizz** to bring your favorite stories to life! 🎉

