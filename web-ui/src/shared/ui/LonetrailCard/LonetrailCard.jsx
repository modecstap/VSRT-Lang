import styles from "./LonetrailCard.module.css";

export default function LonetrailCard({ children }) {
  return (
    <div className={styles.lonetrailWrapper}>
      <div className={styles.lonetrailTopOberation} />

      <div className={styles.lonetrailCard}>
        <div className={styles.lonetrailTricolor} />

        {children}
      </div>

      <div className={styles.lonetrailBottomOberation} />
    </div>
  );
}