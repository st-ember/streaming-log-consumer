# Streaming Log Consumer

A high-performance, scalable companion service to the **StreamingAPI**. This consumer is designed to process log streams independently, allowing for granular scaling and robust error handling.

## Core Features
- **Independent Scaling:** Decoupled from the main API to handle ingestion spikes.
- **Persistent Retry Logic:** Failed log storage attempts are tracked as `LogRetry` entities in a persistent store.
- **Tiered Error Handling:** 
  1. Primary storage attempt.
  2. Automatic retry (up to 5 attempts).
  3. Discarded logs are moved to a dead-letter state for manual investigation.

## Log Format
The consumer expects logs in a structured key-value format:
`category={category} level={level} message={message} id={id} source_id={source_id}`

## Reliability Design
- **Success Management:** Successfully processed logs are cleared from the retry store via a background cron job to keep the storage "hot" and lean.
- **Backoff Strategy:** Calculated by the adapter and stored in the domain to prevent service-hammering during outages.
