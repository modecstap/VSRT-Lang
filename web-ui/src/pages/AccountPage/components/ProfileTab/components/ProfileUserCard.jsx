import { useRef } from "react";
import Button from "../../../../../shared/ui/Button/Button";
import styles from "./ProfileUserCard.module.css";

function ProfileUserCard({ user, uploading, uploadError, handleAvatarChange }) {
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
                    <h1>{uploadError || user.username}</h1>
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
