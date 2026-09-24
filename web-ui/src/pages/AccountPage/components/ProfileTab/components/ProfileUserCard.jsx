import { useRef } from "react";
import Button from "../../../../../shared/ui/Button/Button";
import Input from "../../../../../shared/ui/Input/Input";
import styles from "./ProfileUserCard.module.css";

function ProfileUserCard({
    user,
    uploading,
    uploadError,
    handleAvatarChange,
    editing,
    draft,
    onDraftChange,
    buttonLabel,
    buttonDisabled,
    onProfileButton,
}) {
    const fileInputRef = useRef(null);

    return (
        <div className={styles.user}>
            <div className={styles.avatarWrap}>
                <button
                    type="button"
                    className={styles.avatarButton}
                    aria-label="Upload avatar"
                    disabled={uploading}
                    onClick={() => fileInputRef.current?.click()}
                >
                    <img
                        src={user.avatar}
                        alt={user.username}
                        className={styles.avatar}
                    />
                </button>
                <input
                    ref={fileInputRef}
                    type="file"
                    accept="image/jpeg,image/png,image/webp"
                    className={styles.fileInput}
                    onChange={handleAvatarChange}
                />
            </div>

            <div className={styles.info}>
                <div>
                    {editing ? (
                        <>
                            <Input
                                id="profile-username"
                                name="username"
                                autoComplete="username"
                                value={draft.username}
                                onChange={onDraftChange}
                            />
                            <Input
                                id="profile-email"
                                name="email"
                                type="email"
                                autoComplete="email"
                                value={draft.email}
                                onChange={onDraftChange}
                            />
                        </>
                    ) : (
                        <>
                            <h1>{uploadError || user.username}</h1>
                            <p>{user.email}</p>
                        </>
                    )}
                </div>

                <Button
                    type="button"
                    className={styles.button}
                    disabled={buttonDisabled}
                    onClick={onProfileButton}
                >
                    {buttonLabel}
                </Button>
            </div>
        </div>
    );
}

export default ProfileUserCard;
