import useProfileTab from './hooks/useProfileTab';
import useProfileCard from './hooks/useProfileCard';
import styles from "./ProfileTab.module.css";
import ProfileUserCard from './components/ProfileUserCard';
import SessionList from './components/SessionList';

function ProfileTab() {
    const { uploading, uploadError, handleAvatarChange, ...user } = useProfileCard();
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
                <ProfileUserCard
                    user={user}
                    uploading={uploading}
                    uploadError={uploadError}
                    handleAvatarChange={handleAvatarChange}
                />
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