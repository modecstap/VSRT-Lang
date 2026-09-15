import useProfileTab from './hooks/useProfileTab';
import styles from "./ProfileTab.module.css";
import ProfileUserCard from './components/ProfileUserCard';
import SessionList from './components/SessionList';

function ProfileTab() {
    const user = {
        login: "MODECSTAP",
        email: "TONI.SHEBANIN@MAIL.RU",
        avatar: "https://i.pravatar.cc/150?img=12",
    };

    const {
        sessions,
        loading,
        error,
        creating,
        sessionName,
        handleSessionNameChange,
        handleCreateSession,
        handleSelectSession,
    } = useProfileTab();

    return (
        <div className={styles.container}>
            <div className={styles.card}>
                <ProfileUserCard user={user} />
            </div>

            <SessionList
                loading={loading}
                error={error}
                sessions={sessions}
                sessionName={sessionName}
                creating={creating}
                handleSessionNameChange={handleSessionNameChange}
                handleCreateSession={handleCreateSession}
                handleSelectSession={handleSelectSession}
            />
        </div>
    );
}

export default ProfileTab;