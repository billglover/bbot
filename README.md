# BuddyBot
![Is_a Slack_Bot](https://img.shields.io/badge/Is_a-Slack_Bot-ab47bc.svg?style=flat-square)  ![Uses Serverless_Framework](https://img.shields.io/badge/Uses-Serverless_Framework-brightgreen.svg?style=flat-square)   ![Runs_on Amazon_AWS_Lambda](https://img.shields.io/badge/Runs_on-Amazon_AWS_Lambda-ffad00.svg?style=flat-square)

BuddyBot is a community minded Slack Bot. It allows people to flag Slack messages for possible **Code of Conduct violation**. When a message is flagged three things happen:

+ The user who flagged the message is notified that the report is being looked at.
+ The user who authored the message that has been flagged is notified and asked to review their message.
+ The team admins channel is notified that a message has been flagged, providing details of the message, the name of the reporter and a link to the message.


## Note
The bot has been built using the [**Serverless**](https://serverless.com) Framework & is designed to run on **Amazon AWS Lambda**.
The BuddyBoy uses access tokens.
### **No messages, user data is stored**.

It currently only supports the CodeBuddies Slack workspace but multi workspace support is on the roadmap.