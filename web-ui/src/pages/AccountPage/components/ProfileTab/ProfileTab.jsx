import useProfileTab from './hooks/useProfileTab';
import styles from "./ProfileTab.module.css";
import ProfileUserCard from './components/ProfileUserCard';
import SessionList from './components/SessionList';

const user = {
    login: "MODECSTAP",
    email: "TONI.SHEBANIN@MAIL.RU",
    avatar: "https://i.pravatar.cc/150?img=12",
};

function ProfileTab() {
    const {
        sessions,
        loading,
        error,
        creating,
        deletingSessionId,
        sessionName,
        handleSessionNameChange,
        handleCreateSession,
        handleSelectSession,
        handleDeleteSession,
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
                deletingSessionId={deletingSessionId}
                handleSessionNameChange={handleSessionNameChange}
                handleCreateSession={handleCreateSession}
                handleSelectSession={handleSelectSession}
                handleDeleteSession={handleDeleteSession}
            />
        </div>
    );
}

export default ProfileTab;