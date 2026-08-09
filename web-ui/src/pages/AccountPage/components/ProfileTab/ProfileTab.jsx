import Button from "../../../../shared/ui/Button/Button";
import useProfileTab from './hooks/useProfileTab';
import styles from "./ProfileTab.module.css";

function ProfileTab() {
    const user = {
        login: "MODECSTAP",
        email: "TONI.SHEBANIN@MAIL.RU",
        avatar: "https://i.pravatar.cc/150?img=12",
    };

    const { sessions, loading, error } = useProfileTab();

    return (
        <div className={styles.container}>
            <div className={styles.card}>
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
            </div>

            <section className={styles.card}>
                <h1 className={styles.title}>SESSIONS</h1>

                {loading ? (
                    <p>Loading sessions...</p>
                ) : error ? (
                    <p>{error}</p>
                ) : (
                    <table className={styles.table}>
                        <thead>
                            <tr>
                                <th>DATE</th>
                                <th>NAME</th>
                                <th>SAVED COUNT</th>
                            </tr>
                        </thead>

                        <tbody>
                            {sessions.map((session) => (
                                <tr key={session.id}>
                                    <td>{session.date}</td>
                                    <td>{session.name}</td>
                                    <td>{session.saved}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                )}
            </section>
        </div>
    );
}

export default ProfileTab;