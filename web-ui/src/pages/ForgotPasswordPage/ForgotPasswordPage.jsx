import { Link } from 'react-router-dom';
import Button from '../../shared/ui/Button/Button';
import Input from '../../shared/ui/Input/Input';
import pageStyles from '../LoginPage/LoginPage.module.css';
import formStyles from '../LoginPage/components/LoginForm/LoginForm.module.css';
import { useForgotPasswordForm } from './hooks/useForgotPasswordForm';

function ForgotPasswordPage() {
  const { form, phase, buttonLabel, handleChange, handleSubmit } = useForgotPasswordForm();

  return (
    <main className={pageStyles.page}>
      <section className={pageStyles.panel}>
        <h1 className={pageStyles.title}>FORGOT PASSWORD</h1>
        <form className={formStyles.form} onSubmit={handleSubmit}>
          <div className={formStyles.grid}>
            <Input
              id="forgot-email"
              name="email"
              type="text"
              placeholder="Email"
              value={form.email}
              onChange={handleChange}
              autoComplete="email"
            />
          </div>
          <div className={formStyles.actions}>
            <Button type="submit" disabled={phase === 'submitting'}>
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

export default ForgotPasswordPage;
