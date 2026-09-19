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
  className,
  multiline = false,
  rows = 3,
}) {
  const inputClassName = [styles.input, multiline ? styles.textarea : null, className]
    .filter(Boolean)
    .join(' ');

  return (
    <label className={styles.field} htmlFor={id}>
      {label ? <span className={styles.label}>{label}</span> : null}
      {multiline ? (
        <textarea
          id={id}
          name={name}
          rows={rows}
          className={inputClassName}
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
          className={inputClassName}
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
