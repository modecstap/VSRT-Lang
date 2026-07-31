import { NavLink, Outlet } from 'react-router-dom';
import styles from './AccountPage.module.css';

function AccountPage() {
    return (
        <main className={styles.container}>
            <aside className={styles.sidebar}>
                <NavLink
                    to="profile"
                    className={({ isActive }) =>
                        isActive ? styles.active : styles.link
                    }
                >
                    PROFILE
                </NavLink>

                <NavLink
                    to="session"
                    className={({ isActive }) =>
                        isActive ? styles.active : styles.link
                    }
                >
                    SESSION
                </NavLink>
            </aside>

            <section className={styles.content}>
                <Outlet />
            </section>
        </main>
    );
}

export default AccountPage;