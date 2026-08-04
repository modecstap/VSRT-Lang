import styles from './Input.module.css';

function Input({
  label,
  id,
  name,
  type = 'text',
  placeholder,
  value,
  onChange,
  autoComplete,
  multiline = false,
  rows = 3,
}) {
  return (
    <label className={styles.field} htmlFor={id}>
      {label ? <span className={styles.label}>{label}</span> : null}
      {multiline ? (
        <textarea
          id={id}
          name={name}
          rows={rows}
          className={`${styles.input} ${styles.textarea}`}
          placeholder={placeholder}
          value={value}
          onChange={onChange}
          autoComplete={autoComplete}
        />
      ) : (
        <input
          id={id}
          name={name}
          type={type}
          className={styles.input}
          placeholder={placeholder}
          value={value}
          onChange={onChange}
          autoComplete={autoComplete}
        />
      )}
    </label>
  );
}

export default Input;
