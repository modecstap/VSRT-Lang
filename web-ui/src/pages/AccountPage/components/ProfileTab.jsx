import styles from "./ProfileTab.module.css";
import LonetrailCard from "../../../shared/ui/LonetrailCard/LonetrailCard"

function ProfileTab() {
    const user = {
        login: "MODECSTAP",
        email: "TONI.SHEBANIN@MAIL.RU",
        avatar: "https://i.pravatar.cc/150?img=12",
    };

    const sessions = [
        {
            id: 1,
            date: "01.08.2026",
            name: "Summer project",
            saved: 24,
        },{
            id: 1,
            date: "01.08.2026",
            name: "Summer project",
            saved: 24,
        },
        {
            id: 2,
            date: "29.07.2026",
            name: "Test session",
            saved: 8,
        },
        {
            id: 3,
            date: "22.07.2026",
            name: "Portfolio",
            saved: 41,
        },
    ];

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

                        <button className={styles.button}>
                            CHANGE
                        </button>
                    </div>
                </div>
            </div>

            <section className={styles.card}>
                <h1 className={styles.title}>SESSIONS</h1>

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
            </section>
        </div>
    );
}

export default ProfileTab;