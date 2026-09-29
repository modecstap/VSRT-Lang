import { useState } from 'react';
import styles from './LoginPage.module.css';
import LoginForm from './components/LoginForm/LoginForm';
import RegisterForm from './components/RegisterForm/RegisterForm';
import { useDocumentTitle } from '../../shared/hooks/useDocumentTitle';

function LoginPage() {
  const [isRegistering, setIsRegistering] = useState(false);
  useDocumentTitle(isRegistering ? 'Register' : 'Login');

  const handleSwitchToLogin = () => {
    setIsRegistering(false);
  };

  const handleSwitchToRegister = () => {
    setIsRegistering(true);
  };

  return (
    <main className={styles.page}>
      <section className={styles.panel}>
        <h1 className={styles.title}>{isRegistering ? 'REGISTER' : 'LOGIN'}</h1>
        {isRegistering ? (
          <RegisterForm onSwitchToLogin={handleSwitchToLogin} />
        ) : (
          <LoginForm onSwitchToRegister={handleSwitchToRegister} />
        )}
      </section>
    </main>
  );
}

export default LoginPage;
