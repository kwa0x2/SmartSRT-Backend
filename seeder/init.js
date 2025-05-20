db = db.getSiblingDB('autosrt');

db.createCollection("users");
db.createCollection("usage");
db.createCollection("customer");
db.createCollection("subscription");

db.users.createIndex(
    { email: 1 },
    { 
        unique: true,
        partialFilterExpression: { deleted_at: { $exists: false } }
    }
);

db.users.createIndex(
    { phone_number: 1 },
    { 
        unique: true,
        partialFilterExpression: { deleted_at: { $exists: false } }
    }
);

db.usage.createIndex(
    { user_id: 1 },
    { 
        unique: true,
        partialFilterExpression: { deleted_at: { $exists: false } }
    }
);

db.customer.createIndex(
    {customer_id: 1},
    {
        unique: true,
        partialFilterExpression: { deleted_at: { $exists: false } }
    }
)

db.subscription.createIndex(
    {subscription_id: 1},
    {
        unique: true,
        partialFilterExpression: { deleted_at: { $exists: false } }
    }
)

print("seeder success")