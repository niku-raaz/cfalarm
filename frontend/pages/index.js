import Navbar from '../components/Navbar';

export default function Home() {
  return (
    <div>
      <Navbar />
      <div style={{ padding: '2rem' }}>
        <h2>Welcome to CFAlarm!</h2>
        <p>Automatically track contests and sync them with your Google Calendar.</p>
      </div>
    </div>
  );
}
