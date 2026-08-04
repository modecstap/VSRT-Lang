import styles from './Button.module.css';

function Button({ type = 'button', children, variant = 'primary', className, ...props }) {
  const combinedClassName = [
    styles.button,
    variant === 'primary' ? styles.primary : null,
    className,
  ].filter(Boolean).join(' ');

  return (
    <button type={type} className={combinedClassName} {...props}>
      {children}
    </button>
  );
}

export default Button;
