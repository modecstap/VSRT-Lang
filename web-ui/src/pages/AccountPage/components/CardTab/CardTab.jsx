import Button from '../../../../shared/ui/Button/Button';
import useCardTab from './hooks/useCardTab';
import { ESTIMATIONS } from './model/cardTabModel';
import styles from './CardTab.module.css';
import { useDocumentTitle } from '../../../../shared/hooks/useDocumentTitle';

function CardTab() {
    useDocumentTitle('Cards');
    const {
        canEstimate,
        card,
        estimation,
        handleEstimate,
        handleMain,
        mainDisabled,
        mainLabel,
        passed,
        remaining,
        revealed,
        status,
    } = useCardTab();

    const showBack = Boolean(card) && revealed;

    return (
        <div className={styles.container}>
            <div className={styles.layout}>
                <div className={styles.front}>
                    {status ? (
                        <p className={`${styles.value} ${styles.statusPanel}`}>{status}</p>
                    ) : (
                        <>
                            <p className={`${styles.value} ${styles.phrase}`}>{card.phrase}</p>
                            <p className={`${styles.value} ${styles.baseForm}`}>{card.baseForm}</p>
                            <p className={`${styles.value} ${styles.synonyms}`}>{card.synonyms}</p>
                            <p className={`${styles.value} ${styles.contexts}`}>{card.contexts}</p>
                        </>
                    )}
                </div>
                <div className={styles.back}>
                    <p className={styles.value}>{showBack ? card.translations : ''}</p>
                    <p className={styles.value}>{showBack ? card.contextTranslations : ''}</p>
                </div>
                <div className={styles.menu}>
                    <div className={`${styles.value} ${styles.stats}`}>
                        <p>Remaining: {remaining}</p>
                        <p>Passed: {passed}</p>
                    </div>
                    <div className={styles.estimations}>
                        {ESTIMATIONS.map((item) => (
                            <Button
                                key={item}
                                type="button"
                                className={`${styles.estimation} ${estimation === item ? styles.estimationActive : ''}`}
                                disabled={!canEstimate}
                                onClick={() => handleEstimate(item)}
                            >
                                {item.toLowerCase()}
                            </Button>
                        ))}
                    </div>
                    <Button type="button" className={styles.button} disabled={mainDisabled} onClick={handleMain}>
                        {mainLabel}
                    </Button>
                </div>
            </div>
        </div>
    );
}

export default CardTab;
