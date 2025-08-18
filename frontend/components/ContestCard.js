import React from 'react';


export default function ContestCard({ contest }) {
  return (
    <div className="contest-card">
      <h3>{contest.name}</h3>
      <p>Start: {new Date(contest.startTime).toLocaleString()}</p>
      <p>End: {new Date(contest.endTime).toLocaleString()}</p>
    </div>
  );
}
