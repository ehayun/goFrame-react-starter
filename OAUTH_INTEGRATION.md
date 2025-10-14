# OAuth Integration with Existing Database

## Problem
The initial implementation created a new `users` table schema that conflicted with your existing production database.

## Solution
I've updated the code to work with your **actual** database schema by:

### 1. Updated User Model (`internal/model/user.go`)
Now matches your real database schema with:
- `zehut` as primary key (Israeli ID)
- `first_name` + `last_name` (not just `name`)
- `hashed_password` for existing authentication
- All 47 fields from your actual users table
- Proper field mappings and JSON serialization

### 2. Updated Repository (`internal/repository/user_repository.go`)
- `FindByZehut()` - Find by primary key (zehut)
- `FindByEmail()` - Find by email for OAuth matching
- Removed non-existent fields like `google_id`, `deleted_at`
- Uses `inserted_at` and `updated_at` (not `created_at`)

### 3. OAuth Flow with Existing Users
**How it works:**
1. User logs in with Google OAuth
2. System receives email from Google
3. **Looks up user by email in existing users table**
4. If found → Creates session and logs them in
5. If not found → Shows "not authorized" message
6. Updates `avatar` field with Google profile picture
7. Sets `confirmed_at` timestamp on first OAuth login

**Key Points:**
- Only existing users (already in your database) can log in via OAuth
- Matches users by **email** field
- No new user creation (prevents unauthorized access)
- Preserves all existing user data and permissions

### 4. Session Management
Session now stores:
- `zehut` - Primary key from users table
- `email` - User's email
- `name` - Full name (first_name + last_name)

### 5. Removed Migration Files
Deleted incorrect migration files that would have created a conflicting users table.

## What OAuth Adds to Existing Users
When a user logs in via Google OAuth:
- **Updates `avatar`** with Google profile picture
- **Sets `confirmed_at`** on first OAuth login
- **Creates Redis session** for authentication
- **No other data is modified**

## Database Changes Required
**NONE** - The code now works with your existing schema as-is.

Optional: If you want to track which users use OAuth, you could add:
```sql
ALTER TABLE users ADD COLUMN google_oauth_enabled BOOLEAN DEFAULT FALSE;
```

## Authentication Flow
```
User → Google OAuth → Email from Google
           ↓
    Find user by email in existing DB
           ↓
        Found?
      /         \
    YES          NO
     ↓           ↓
Create session   Show error
     ↓           (not authorized)
 Log user in
```

## Testing
To test with your database:
1. Ensure a user exists with a valid email
2. Use that email to log in via Google OAuth
3. System will match and authenticate the existing user
4. Check `avatar` and `confirmed_at` fields are updated

## Security
- Only pre-existing users can log in
- No automatic user creation
- Maintains existing role-based access control
- Session managed in Redis with TTL
