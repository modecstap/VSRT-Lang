import styles from './LoginForm.module.css';
import Button from '../../../../shared/ui/Button/Button';
import FormField from './FormField';
import { useLoginForm } from '../../hooks/useLoginForm';

function LoginForm({ onSwitchToRegister }) {
  const { form, error, isSubmitting, handleChange, handleSubmit } = useLoginForm();

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.grid}>
        <FormField
          label=" "
          id="login-username"
          name="login"
          type="text"
          placeholder="Login"
          value={form.login}
          onChange={handleChange}
        />

        <FormField
          label=" "
          id="login-password"
          name="password"
          type="password"
          placeholder="Password"
          value={form.password}
          onChange={handleChange}
          autoComplete="current-password"
        />
      </div>

      <div className={styles.actions}>
        <Button type="submit" disabled={isSubmitting}>
          {error ? 
            error : 
            isSubmitting ? 'Signing In...' : 'Sign In'
          }
          
        </Button>
        <div className={styles.additional}>
          <a className={styles.link} href="/forgot-password">
            FORGOT PASSWORD
          </a>
          <button type="button" className={styles.link} onClick={onSwitchToRegister}>
            CREATE ACCOUNT
          </button>
        </div>
      </div>
    </form>
  );
}

export default LoginForm;
