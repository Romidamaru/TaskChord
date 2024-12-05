# TaskChord

TaskChord is a lightweight Discord bot designed to help users manage tasks within Discord servers. Each server (guild) acts as a namespace, isolating tasks to ensure that your to-dos stay organized and relevant to the respective community.

Features

	•	Create Tasks: Easily add tasks with a title, description, and optional priority (High, Medium, Low).
	•	Show Tasks: View all tasks or search for a specific task by ID.
	•	Delete Tasks: Remove tasks by specifying their ID.
	•	Namespace Isolation: Tasks are isolated per Discord server, ensuring privacy and organization.

### Commands

1. /create

Creates a new task.

Options:
•	title (Required): The title of the task.
•	description (Required): A detailed description of the task.
•	priority (Optional): The priority of the task (High, Medium, Low).

Example:

/create title: "Buy groceries" description: "Milk, eggs, bread" priority: "High"

Response:

	Task #1 “Buy groceries” successfully created!

2. /show

Displays your tasks.

Options:
•	id (Optional): The ID of a specific task to view.

Example:
•	Show all tasks:

/show


	•	Show a specific task:

/show id: "1"



Response:
If tasks exist:

Your Tasks:
- #1 "Buy groceries"
  Priority: High
  Description: Milk, eggs, bread

If no tasks exist:

You have no tasks!

3. /delete

Deletes a task by ID.

Options:
•	id (Required): The ID of the task to delete.

Example:

/delete id: "1"

Response:

	Task #1 successfully deleted!

### Setup

	1.	Clone the Repository:

git clone <repository-url>
cd TaskChord


	2.	Install Dependencies:

go mod tidy


	3.	Configure Environment Variables:
	•	Create a .env file in the project root.
	•	Add your Discord bot token and database configuration:

DISCORD_BOT_TOKEN=<your-bot-token>
DATABASE_URL=<your-database-url>


	4.	Run the Bot:

go run cmd/bot/main.go

### Database Schema

The bot uses a database to store tasks. Each task is associated with a guild_id (Discord server) and a user_id (task owner). The schema includes:
•	task_id_in_guild: A unique ID for tasks within each server.
•	guild_id: Discord server ID.
•	user_id: Discord user ID.
•	title: Task title.
•	description: Task description.
•	priority: Task priority (High, Medium, Low).

### Future Enhancements

	•	Add due dates to tasks.
	•	Enable task updates.
	•	Add categories or tags for better task organization.
	•	Integrate reminders for upcoming tasks.

### Contributing

We welcome contributions! Please feel free to submit issues or pull requests to improve TaskChord.

### License

This project is licensed under the MIT License. See the LICENSE file for details.

Enjoy organizing your tasks with TaskChord! 🎉