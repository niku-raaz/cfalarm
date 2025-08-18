# CFAlarm: Automated Contest Tracker and Practice Companion

CFAlarm is a web application designed for competitive programmers to automate their contest participation and streamline their practice routines. It automatically registers users for upcoming Codeforces contests, adds them to their Google Calendar, and provides daily practice problems tailored to their rating.

---

## Core Features

- **Google Authentication**: Secure and easy sign-in using your Google account.
- **User Profile Management**: Update your competitive programming platform handles (Codeforces, Codechef, Atcoder, Leetcode).
- **Automated Contest Registration**: A scheduled service runs every two days to find upcoming Codeforces contests and automatically register you.
- **Google Calendar Integration**: Automatically creates events in your primary Google Calendar for every contest you are registered for.
- **Daily Practice Problems**: Receive two new unsolved problems from Codeforces every day at midnight, scaled to your current rating (x+100 and x+200).
- **Practice Tracking**: View your recent unsolved practice problems and mark them as solved.
- **Email Reminders**: Get an email notification at 9 PM if you have not completed your daily practice problems.

---

## Tech Stack

### Backend
- **Language**: Go
- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT and Google OAuth2
- **Scheduling**: `robfig/cron` for cron jobs
- **Email Service**: SendGrid

### Frontend
- **Framework**: Next.js
- **Library**: React
- **Styling**: Standard CSS
- **API Communication**: Axios



