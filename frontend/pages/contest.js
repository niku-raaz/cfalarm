import { useEffect, useState } from 'react';
import Navbar from '../components/Navbar';
import ContestCard from '../components/ContestCard';
import axios from 'axios';

export default function Contests() {
  const [contests, setContests] = useState([]);

  useEffect(() => {
    axios.get('http://localhost:8080/api/contests') // replace with your backend contests endpoint
      .then(res => setContests(res.data))
      .catch(err => console.log(err));
  }, []);

  return (
    <div>
      <Navbar />
      <div style={{ padding: '2rem' }}>
        <h2>Upcoming Contests</h2>
        <div>
          {contests.map(contest => (
            <ContestCard key={contest.id} contest={contest} />
          ))}
        </div>
      </div>
    </div>
  );
}
