import { useState } from 'react';
import styles from './LoginPage.module.css';
import LoginForm from './components/LoginForm/LoginForm';
import RegisterForm from './components/RegisterForm/RegisterForm';
import LonetrailCard from '../../shared/ui/LonetrailCard/LonetrailCard';

function LoginPage() {
  const [isRegistering, setIsRegistering] = useState(false);

  const handleSwitchToLogin = () => {
    setIsRegistering(false);
  };

  return (
    <main className={styles.page}>
      <div>
        <section className={styles.panel}>
          <h1 className={styles.title}>{isRegistering ? 'REGISTER' : 'LOGIN'}</h1>
          {isRegistering ? (
            <RegisterForm onSwitchToLogin={handleSwitchToLogin} />
          ) : (
            <LoginForm onSwitchToRegister={() => setIsRegistering(true)} />
          )}
        </section>
      </div>
    </main>
  );
}

export default LoginPage;
