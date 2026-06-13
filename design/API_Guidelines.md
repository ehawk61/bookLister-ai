# API Design Guidelines 

## Background 
In thinking on this API, I wanted to capture up-front the general guidelines for the API before encapuslating them into the application or spec in order to provide a human-readable 'why' to the choices made. 

Generally speaking, my goals are as follows for the API: 
- Concise Usage of HTTP Codes to provide other clients of the API a way to 'reason' on what actions to take
- A Human-readable HTTP body in errors in an effort to help a developer understand what has transpired without revealing too much that becomes a security risk 

## HTTP Codes to Be Used
**Note: Some of these choices were suggested in a discussion with Claude Code. As the main developer, I agree with the suggestions provided by Claude Code but reconigze the possiblity that what has been suggested might be incorrect.** 
| HTTP Code | Allowed Verbs | Meaning                                       |
| --------- | ------------- | --------------------------------------------- |
| 200       | GET, HEAD     | Retrieved an Entity                           |
| 201       | POST          | Created a New Entity                          |
| 204       | PUT, DELETE   | Successfully Updated or Deleted the Entity    |
| 207       | POST, PUT     | Multiple Entities provided with Mixed Results |
| 422       | POST, PUT     | The Entity Fails Validation                   |
| 400       | ALL           | Bad Request                                   |
| 404       | GET           | Not Found                                     |

## Error Response Design
All errors should have three parts to them: 
- `error_code`
  - A concise machine-readable string that provides a client developer the ability to "reason" next steps
- `error_description`
  - A concise human-readable string that provide an operations or client developer
- `details`
  - Details on what has gone wrong with the caveat of limiting details that should only be visible to the developers 
    - Anti-Example: Dumping straight stack trace into this field 

### Error Response Setup

| HTTP Code | Error Code              | Error Description                                   | Details                  |
| --------- | ----------------------- | --------------------------------------------------- | ------------------------ |
| 404       | BOOK_NOT_FOUND          | Book does not exist                                 |                          |
| 404       | AUTHOR_NOT_FOUND        | Author does not exist                               |                          |
| 422       | BOOK_VALIDATION_ERROR   | The book request is not valid                       | Put in the offending bit |
| 422       | AUTHOR_VALIDATION_ERROR | The author request is not valid                     | Put in the offending bit |
| 409       | DUPLICATE_BOOK_ENTRY    | This book entry matches too closely another entry   | Provide the close match  |
| 409       | DUPLICATE_AUTHOR_ENTRY  | This author entry matches too closely another entry | Provide the close match  |

## API Structure

### bookLister Headers 
Any API Request to the bookLister API will be required to contain the following headers:
- X-API-Version (Optional)
   - Indicates which version of the bookLister API is to be used to ensure compatability 
   - If no header is provided then a fallback version will be used   
- X-Correlation-ID (Required)
   - Allows for tracability in logs 
   - Expected to be a UUIDv4 entry   

### bookLister Endpoints
Over time this list will change but its a best guess at the moment 

| HTTP Verb | URI                | URI Intention                                                | Development Notes    |
| --------- | ------------------ | ------------------------------------------------------------ | -------------------- |
| GET       | /books             | Get all books within the database                            | pagination is needed |
| POST      | /books             | Add books to the database                                    | 1 or many books      |
| GET       | /books/{id}        | Get a single book by its id                                  |                      |
| PUT       | /books/{id}        | Update a book by its id                                      |                      |
| DELETE    | /books/{id}        | Delete a book by its id                                      |                      |
| GET       | /books/q?title=""  | Search for all books by title that contains provided string  | pagination is needed |
| GET       | /books/q?author="" | Search for all books by author that contains provided string | pagination is needed |
| GET       | /authors           | Get all the authors within the database                      | pagination is needed | 
| POST      | /authors           | Add author(s) to the database                                | 1 or many authors    |
| GET       | /authors/{id}      | Get a single author by its id                                |                      |
| PUT       | /authors/{id}      | Update an author by its id                                   |                      |
| DELETE    | /authors/{id}      | Delete an author by its id                                   |                      |
| GET       | /authors/q?name="" | Search for all authors by name that contains provided string | pagination is needed | 
| GET       | /healthcheck       | Check the status of the bookLister API                       |                      |
