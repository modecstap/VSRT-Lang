import styles from './Button.module.css';

function Button({ type = 'button', children, variant = 'primary', ...props }) {
  const className =
    variant === 'primary' ? `${styles.button} ${styles.primary}` : styles.button;

  return (
    <button type={type} className={className} {...props}>
      {children}
    </button>
  );
}

export default Button;
