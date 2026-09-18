import Input from "../../../../../shared/ui/Input/Input";
import Button from "../../../../../shared/ui/Button/Button";
import styles from "./SessionList.module.css";

function SessionList({
    loading,
    error,
    sessions,
    sessionName,
    creating,
    deletingSessionId,
    handleSessionNameChange,
    handleCreateSession,
    handleSelectSession,
    handleDeleteSession,
}) {
    return (
        <section className={styles.card}>
            <div className={styles.header}>
                <h1 className={styles.title}>SESSIONS</h1>
                <div className={styles.createSession}>
                    <Input
                        id="sessionName"
                        name="sessionName"
                        placeholder="Session name"
                        value={sessionName}
                        onChange={handleSessionNameChange}
                    />
                    <Button
                        type="button"
                        className={styles.newButton}
                        onClick={handleCreateSession}
                        disabled={creating}
                    >
                        NEW
                    </Button>
                </div>
            </div>

            {loading ? (
                <table className={styles.table}>
                    <thead>
                        <tr>
                            <th>Loading sessions...</th>
                        </tr>
                    </thead>
                </table>
            ) : error ? (
                <table className={styles.table}>
                    <thead>
                        <tr>
                            <th>{error}</th>
                        </tr>
                    </thead>
                </table>
            ) : (
                <table className={styles.table}>
                    <thead>
                        <tr>
                            <th>DATE</th>
                            <th>NAME</th>
                            <th>SAVED COUNT</th>
                            <th></th>
                        </tr>
                    </thead>

                    <tbody>
                        {sessions.map((session) => (
                            <tr
                                key={session.id}
                                onClick={() => handleSelectSession(session)}
                                onKeyDown={(event) => {
                                    if (event.key === 'Enter' || event.key === ' ') {
                                        event.preventDefault();
                                        handleSelectSession(session);
                                    }
                                }}
                                role="button"
                                tabIndex={0}
                            >
                                <td>{session.date}</td>
                                <td>{session.name}</td>
                                <td>{session.saved}</td>
                                <td>
                                    <Button
                                        type="button"
                                        className={styles.deleteButton}
                                        onClick={(event) => {
                                            event.stopPropagation();
                                            handleDeleteSession(session);
                                        }}
                                        disabled={deletingSessionId === session.id}
                                    >
                                        {deletingSessionId === session.id ? '0' : 'X'}
                                    </Button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </section>
    );
}

export default SessionList;
