import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import { useSWRConfig } from 'swr';
import { clearAuthTokens } from '../../api/auth';
import { clearActiveSessionId } from '../../api/sessions';
import { prefetchSessionTab } from '../../routes/prefetch';
import styles from './AccountPage.module.css';

function AccountPage() {
    const navigate = useNavigate();
    const { mutate } = useSWRConfig();

    const handleLogout = () => {
        clearAuthTokens();
        clearActiveSessionId();
        mutate(() => true, undefined, { revalidate: false });
        navigate('/');
    };

    return (
        <div className={styles.page}>
            <main className={styles.container}>
                <aside className={styles.sidebar}>
                    <h1 className={styles.title}>MENU</h1>

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
                        onFocus={prefetchSessionTab}
                        onMouseEnter={prefetchSessionTab}
                    >
                        SESSION
                    </NavLink>

                    <button type="button" className={styles.link} onClick={handleLogout}>
                        LOGOUT
                    </button>
                </aside>

                <section className={styles.content}>
                    <Outlet />
                </section>
            </main>

            <footer className={styles.footer}></footer>
        </div>
    );
}

export default AccountPage;
