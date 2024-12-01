<h1>OCR-React-native-backend</h1>
<hr><p>This project implements a backend system that provides the following functionality:</p>
<p>User Authentication: Sign-up, login, and secure user sessions.</p>
<p>Image-to-Text Conversion: Extract text from uploaded images using Tesseract OCR.</p>
<p>Data Storage and Retrieval: Save extracted text and associated metadata (e.g., image name) in a database and allow users to retrieve their extracted data.</p>
<p>The API is built using Go with Gin for HTTP routing, github.com/tiagomelo/go-ocr/ocr for text extraction, and a PostgreSQL database for data storage.</p><h2>General Information</h2>
<hr><ul>
<li>User Authentication</li>
</ul>
<p>Users can register and log in using secure endpoints.
Authentication tokens are used to manage user sessions.
Each user's data is securely associated with their account.</p><ul>
<li>Image Upload and OCR</li>
</ul>
<p>Users can upload an image through the API.
The system processes the image using Tesseract OCR and extracts readable text.
The extracted text and metadata are stored in the database.</p><ul>
<li>View Extracted Data</li>
</ul>
<p>Users can retrieve a list of their previously uploaded images and the extracted text.</p><h2>Technologies Used</h2>
<hr><ul>
<li>Go lang</li>
</ul><ul>
<li>Postgres</li>
</ul><ul>
<li>Sqlc</li>
</ul><ul>
<li>Tesseract OCR</li>
</ul><h2>Features</h2>
<hr><ul>
<li>User Authentication</li>
</ul><ul>
<li>Image Upload and OCR</li>
</ul><ul>
<li>View Extracted Data</li>
</ul>
