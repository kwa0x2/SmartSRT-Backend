db = db.getSiblingDB("autosrt");

db.createCollection("users");
db.createCollection("usage");
db.createCollection("customer");
db.createCollection("subscription");

db.users.createIndex(
    { email: 1 },
    {
        unique: true,
        partialFilterExpression: { deleted_at: null },
    }
);

db.users.createIndex(
    { phone_number: 1 },
    {
        unique: true,
        partialFilterExpression: { deleted_at: null },
    }
);

db.usage.createIndex(
    { user_id: 1 },
    {
        unique: true,
        partialFilterExpression: { deleted_at: null },
    }
);

db.customer.createIndex(
    { customer_id: 1 },
    {
        unique: true,
        partialFilterExpression: { deleted_at: null },
    }
);

db.subscription.createIndex(
    { subscription_id: 1 },
    {
        unique: true,
        partialFilterExpression: { deleted_at: null },
    }
);

print("✅ Collections and indexes created.");
