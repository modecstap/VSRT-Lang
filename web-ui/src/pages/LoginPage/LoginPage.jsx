import { useState } from 'react';
import styles from './LoginPage.module.css';
import LoginForm from './components/LoginForm/LoginForm';
import RegisterForm from './components/RegisterForm/RegisterForm';
import LonetrailCard from '../../shared/ui/LonetrailCard/LonetrailCard';

function LoginPage() {
  const [isRegistering, setIsRegistering] = useState(false);

  return (
    <main className={styles.page}>
      <LonetrailCard>
        <section className={styles.panel}>
          <h1 className={styles.title}>{isRegistering ? 'REGISTER' : 'LOGIN'}</h1>
          {isRegistering ? (
            <RegisterForm onSwitchToLogin={() => setIsRegistering(false)} />
          ) : (
            <LoginForm onSwitchToRegister={() => setIsRegistering(true)} />
          )}
        </section>
      </LonetrailCard>
    </main>
  );
}

export default LoginPage;
