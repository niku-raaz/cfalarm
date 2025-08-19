import React from 'react';

export default function ContestCard({ contest }) {
  // The API provides time in seconds, so we multiply by 1000 for JavaScript's milliseconds.
  const startTime = new Date(contest.startTimeSeconds * 1000);

  // Calculate duration in hours
  const durationInHours = contest.durationSeconds / 3600;

  const options = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true,
  };

  const timeInUTC = startTime.toLocaleString('en-US', { ...options, timeZone: 'UTC' });
  const timeInIST = startTime.toLocaleString('en-IN', { ...options, timeZone: 'Asia/Kolkata' });
  
  // Check if the contest is finished to conditionally show the link
  const isFinished = contest.phase === 'FINISHED';
  const contestUrl = `https://codeforces.com/contest/${contest.id}`;

  return (
    <div className="contest-card">
      <h3>{contest.name}</h3>
      <p>Start Time (IST): {timeInIST}</p>
      <p>Start Time (UTC): {timeInUTC}</p>
      <p>Duration: {durationInHours} hrs</p>

      {/* --- ADD THIS SECTION --- */}
      {/* If the contest is finished, display a link to the contest page */}
      {isFinished && (
        <div style={{ marginTop: '1rem', textAlign: 'center' }}>
          <a 
            href={contestUrl} 
            target="_blank" 
            rel="noopener noreferrer" 
            style={{
              color: '#fff',
              backgroundColor: '#007bff',
              padding: '0.5rem 1rem',
              borderRadius: '5px',
              textDecoration: 'none'
            }}
          >
            View Contest
          </a>
        </div>
      )}
      {/* --- END OF NEW SECTION --- */}
    </div>
  );
}