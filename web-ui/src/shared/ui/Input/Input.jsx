import styles from './Input.module.css';

function Input({ label, id, name, type = 'text', placeholder, value, onChange, autoComplete }) {
  return (
    <label className={styles.field} htmlFor={id}>
      {label ? <span className={styles.label}>{label}</span> : null}
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
    </label>
  );
}

export default Input;
