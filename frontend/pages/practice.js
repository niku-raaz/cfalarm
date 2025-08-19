import { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import Navbar from '../components/Navbar';
import axios from 'axios';

export default function Practice() {
  const router = useRouter();
  const [problems, setProblems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [rating, setRating] = useState(1200); // Default rating

  useEffect(() => {
    const fetchProblems = async () => {
      const token = localStorage.getItem('authToken');
      if (!token) {
        router.replace('/login');
        return;
      }

      setLoading(true);
      setError('');
      try {
        const response = await axios.get(`http://localhost:8080/api/practice/problems?rating=${rating}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        setProblems(response.data);
      } catch (err) {
        const errorMessage = err.response?.data?.error || 'Failed to fetch problems. Please try again.';
        setError(errorMessage);
        setProblems([]); // Clear previous problems on error
      } finally {
        setLoading(false);
      }
    };

    fetchProblems();
  }, [rating, router]); // Refetch when rating changes

  const handleRatingChange = (e) => {
    setRating(e.target.value);
  };

  return (
    <div>
      <Navbar />
      <div style={{ padding: '2rem' }}>
        <h2>Practice Problems</h2>
        <p>Find unsolved problems from Codeforces based on your preferred rating.</p>
        
        <div style={{ margin: '2rem 0' }}>
          <label htmlFor="rating-select">Select Problem Rating: </label>
          <select id="rating-select" value={rating} onChange={handleRatingChange}>
            {/* Create options for ratings from 800 to 3500 */}
            {Array.from({ length: 28 }, (_, i) => 800 + i * 100).map(r => (
              <option key={r} value={r}>{r}</option>
            ))}
          </select>
        </div>

        {loading && <p>Loading problems...</p>}
        {error && <p style={{ color: 'red' }}>Error: {error}</p>}
        
        {!loading && !error && (
          <div className="problems-list">
            {problems.length > 0 ? problems.map(problem => (
              <div key={`${problem.contestId}-${problem.index}`} style={problemCardStyle}>
                <a 
                  href={`https://codeforces.com/problemset/problem/${problem.contestId}/${problem.index}`} 
                  target="_blank" 
                  rel="noopener noreferrer"
                  style={{ textDecoration: 'none', color: '#007bff' }}
                >
                  {problem.name}
                </a>
                <p style={{ margin: '0.5rem 0 0', color: '#6c757d' }}>Rating: {problem.rating}</p>
              </div>
            )) : <p>No unsolved problems found for this rating.</p>}
          </div>
        )}
      </div>
    </div>
  );
}

// Basic styling for the problem cards
const problemCardStyle = {
  background: '#fff',
  border: '1px solid #ddd',
  borderRadius: '8px',
  padding: '1rem',
  marginBottom: '1rem',
  boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
};