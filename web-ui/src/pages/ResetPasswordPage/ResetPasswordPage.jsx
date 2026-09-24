import { Link } from 'react-router-dom';
import Button from '../../shared/ui/Button/Button';
import Input from '../../shared/ui/Input/Input';
import pageStyles from '../LoginPage/LoginPage.module.css';
import formStyles from '../LoginPage/components/LoginForm/LoginForm.module.css';
import { useResetPasswordForm } from './hooks/useResetPasswordForm';

function ResetPasswordPage() {
  const { form, phase, buttonLabel, handleChange, handleSubmit } = useResetPasswordForm();

  return (
    <main className={pageStyles.page}>
      <section className={pageStyles.panel}>
        <h1 className={pageStyles.title}>RESET PASSWORD</h1>
        <form className={formStyles.form} onSubmit={handleSubmit}>
          <div className={formStyles.grid}>
            <Input
              id="reset-password"
              name="password"
              type="password"
              placeholder="Password"
              value={form.password}
              onChange={handleChange}
              autoComplete="new-password"
            />
            <Input
              id="reset-repeat-password"
              name="repeatPassword"
              type="password"
              placeholder="Repeat password"
              value={form.repeatPassword}
              onChange={handleChange}
              autoComplete="new-password"
            />
          </div>
          <div className={formStyles.actions}>
            <Button type="submit" disabled={phase === 'submitting' || phase === 'success'}>
              {buttonLabel}
            </Button>
            <div>
              <Link className={formStyles.link} to="/">
                BACK TO LOGIN
              </Link>
            </div>
          </div>
        </form>
      </section>
    </main>
  );
}

export default ResetPasswordPage;
