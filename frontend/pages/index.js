import { useState, useEffect } from 'react';
import Navbar from '../components/Navbar';
import ContestCard from '../components/ContestCard';
import axios from 'axios';

export default function Home() {
  const [upcomingContests, setUpcomingContests] = useState([]);
  const [finishedContests, setFinishedContests] = useState([]); // State for finished contests
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchAllContests = async () => {
      try {
        // Fetch both upcoming and finished contests at the same time
        const [upcomingRes, finishedRes] = await Promise.all([
          axios.get('http://localhost:8080/api/contests/upcoming'),
          axios.get('http://localhost:8080/api/contests/finished')
        ]);
        setUpcomingContests(upcomingRes.data);
        setFinishedContests(finishedRes.data);
      } catch (error) {
        console.error('Failed to fetch contests:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchAllContests();
  }, []);

  return (
    <div>
      <Navbar />
      <div style={{ padding: '2rem' }}>
        <h2>Welcome to CFAlarm!</h2>
        <p>Automatically track contests and sync them with your Google Calendar.</p>
        
        <h3 style={{ marginTop: '3rem' }}>Upcoming Contests</h3>
        {loading ? (
          <p>Loading contests...</p>
        ) : (
          <div className="contests-container">
            {upcomingContests && upcomingContests.length > 0 ? (
              upcomingContests.map(contest => (
                <ContestCard key={contest.id} contest={contest} />
              ))
            ) : (
              <p>No upcoming contests found.</p>
            )}
          </div>
        )}

        {/* --- NEW SECTION FOR PAST CONTESTS --- */}
        <h3 style={{ marginTop: '3rem' }}>Past Contests</h3>
        {loading ? (
          <p>Loading contests...</p>
        ) : (
          <div className="contests-container">
            {finishedContests && finishedContests.length > 0 ? (
              finishedContests.map(contest => (
                <ContestCard key={contest.id} contest={contest} />
              ))
            ) : (
              <p>No past contests found.</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}