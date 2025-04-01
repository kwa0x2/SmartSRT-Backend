db = db.getSiblingDB('autosrt');

db.createCollection("users");
db.createCollection("usage");

db.users.createIndex(
    { email: 1 },
    { unique: true }
);

db.users.createIndex(
    { phone_number: 1 },
    { unique: true }
);

db.usage.createIndex(
    {user_id: 1},
    {unique: true}
);

print("seeder success")