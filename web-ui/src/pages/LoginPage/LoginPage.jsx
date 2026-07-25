import styles from './LoginPage.module.css';
import LoginForm from './components/LoginForm/LoginForm';
import LonetrailCard from "../../shared/ui/LonetrailCard/LonetrailCard";

function LoginPage() {
  return (
    <main className={styles.page}>
      <LonetrailCard>
        <section className={styles.panel}>
          <h1 className={styles.title}>LOGIN</h1>
          <LoginForm />
        </section>
      </LonetrailCard>
    </main>
  );
}

export default LoginPage;
