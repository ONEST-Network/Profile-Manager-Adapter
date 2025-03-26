package redis

import (
    "encoding/json"
    "time"
)


// Set stores a value in Redis with an optional expiration
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
    bytes, err := json.Marshal(value)
    if err != nil {
        return err
    }

    return r.Client.Set(r.Context, key, bytes, expiration).Err()
}

// Get retrieves a value from Redis and unmarshals it into the provided interface
func (r *RedisClient) Get(key string, dest interface{}) error {
    bytes, err := r.Client.Get(r.Context, key).Bytes()
    if err != nil {
        return err
    }

    return json.Unmarshal(bytes, dest)
}

// Delete removes a key from Redis
func (r *RedisClient) Delete(key string) error {
    return r.Client.Del(r.Context, key).Err()
}

// Exists checks if a key exists in Redis
func (r *RedisClient) Exists(key string) (bool, error) {
    result, err := r.Client.Exists(r.Context, key).Result()
    return result > 0, err
}
