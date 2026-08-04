import Button from '../../../../shared/ui/Button/Button';
import Input from '../../../../shared/ui/Input/Input';
import useSessionTab from './hooks/useSessionTab';
import styles from './SessionTab.module.css';

function SessionTab() {
    const {
        details,
        form,
        handleContextChange,
        handleSelectWord,
        handleWordChange,
        handleWrite,
        selectedWord,
        words,
    } = useSessionTab();

    return (
        <div className={styles.container}>
            <div className={styles.layout}>
                <div className={`${styles.panel} ${styles.listCard}`}>
                    <div className={styles.list}>
                        {words.map((word) => (
                            <Button
                                key={word.id}
                                className={`${styles.wordItem} ${selectedWord?.id === word.id ? styles.wordItemActive : ''}`}
                                type="button"
                                onClick={() => handleSelectWord(word)}
                            >
                                <span>- {word.word}</span>
                            </Button>
                        ))}
                    </div>
                </div>
                <div className={`${styles.panel} ${styles.wordCard}`}>
                    {details ? (
                        <>
                            <p className={styles.value}>{details.translation}</p>
                            <p className={styles.value}>{details.synonyms.join(', ') || '—'}</p>
                            <p className={styles.value}>{details.antonyms.join(', ') || '—'}</p>
                        </>
                    ) : (
                        <p className={styles.value}>Select or type a word to see details.</p>
                    )}
                </div>
                <div className={`${styles.panel} ${styles.card}`}>
                    {details ? (
                        <>
                            <p className={styles.value}>{details.meaning}</p>
                            <p className={styles.value}>{details.contexts.join(' / ') || '—'}</p>
                        </>
                    ) : (
                        <p className={styles.value}>Add a new word to populate this section.</p>
                    )}
                </div>
                <div className={`${styles.panel} ${styles.contextField}`}>
                    <Input
                        className={styles.input}
                        id="context"
                        name="context"
                        placeholder="Type a context"
                        value={form.context}
                        onChange={handleContextChange}
                        multiline={true}
                        rows={5}
                    />
                </div>
                <div className={`${styles.panel} ${styles.wordField}`}>
                    <Input
                        className={styles.input}
                        id="word"
                        name="word"
                        placeholder="Type a word"
                        value={form.word}
                        onChange={handleWordChange}
                        multiline={true}
                        rows={5}
                    />
                </div>
                <div className={`${styles.panel} ${styles.submitButton}`}>
                    <Button type="button" className={styles.button} onClick={handleWrite}>WRITE</Button>
                </div>
            </div>           
        </div>
    );
}

export default SessionTab;