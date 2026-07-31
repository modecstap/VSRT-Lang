import styles from "./LonetrailCard.module.css";

export default function LonetrailCard({ children }) {
  return (
    <div className={styles.wrapper}>
      <div className={styles.topDecoration} />

      <div className={styles.card}>
        <div className={styles.accentBar} />

        {children}
      </div>

      <div className={styles.bottomDecoration} />
    </div>
  );
}