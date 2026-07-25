import styles from './LoginForm.module.css';
import Button from '../../../../shared/ui/Button/Button';
import FormField from './FormField';
import { useLoginForm } from '../../hooks/useLoginForm';

function LoginForm() {
  const { form, error, isSubmitting, handleChange, handleSubmit } = useLoginForm();

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.grid}>
        <FormField
          label=" "
          id="login-username"
          name="username"
          type="text"
          placeholder="Username"
          value={form.username}
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
        <a className={styles.link} href="/forgot-password">
          FORGOT PASSWORD
        </a>
      </div>
    </form>
  );
}

export default LoginForm;
