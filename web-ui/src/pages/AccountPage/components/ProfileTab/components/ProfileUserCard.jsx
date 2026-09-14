import Button from "../../../../../shared/ui/Button/Button";
import styles from "./ProfileUserCard.module.css";

function ProfileUserCard({ user }) {
    return (
        <div className={styles.user}>
            <img
                src={user.avatar}
                alt={user.login}
                className={styles.avatar}
            />

            <div className={styles.info}>
                <div>
                    <h1>{user.login}</h1>
                    <p>{user.email}</p>
                </div>

                <Button className={styles.button}>
                    CHANGE
                </Button>
            </div>
        </div>
    );
}

export default ProfileUserCard;
