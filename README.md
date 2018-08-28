# BuddyBot

## Development Log

- Create base serverless project `serverless create -t aws-go -n bbot -p bbot`
- Create `functions` folder and delete the `hello` and `world` functions in the root of the project
- Create functions for our application
  - `acceptRequest` - to handle all inbound web-hook requests from Slack
  - `flagMessage` - to handle flagged messages appropriately
- Modify `Makefile` and `serverless.yml` to point to our new functions
- Create `flagMessageQueue` and use as the event source for `flagMessage`
- Test we can call the HTTP endpoint for `acceptRequest`
- Test we can place messages on `flagMessageQueue`
