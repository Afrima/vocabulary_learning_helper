import React, {useEffect, useMemo, useState} from 'react';
import Grid from "../../components/ui/Grid/Grid";
import {deleteCall, get, post} from "../../utility/restCaller";
import VocabularyEditModal from "../../components/ui/Modal/VocabularyEditModal";

export type Vocab = {
    id?: string;
    german: string;
    japanese: string;
    kanji: string;
}

const VocabularyView = (): JSX.Element => {
    const emptyEditData = {german: '', japanese: '', kanji: ''};
    const [vocabs, setVocabs] = useState<Vocab[]>([]);
    const [editData, setEditData] = useState<Vocab>(emptyEditData);
    const [showEditModal, setShowEditModal] = useState<boolean>(false);


    useEffect(() => {
        get<Vocab[]>('/vocab', setVocabs);
    }, []);

    const grid = useMemo(() => {
        const deleteHandler = (data: Vocab): void => {
            deleteCall<Vocab, Vocab>('/vocab', data, ((d) => setVocabs(vocabs.filter(vocab => vocab.id !== d.id))));
        };
        const setEditHandler = (data: Vocab): void => {
            setEditData(data);
            setShowEditModal(true);
        };
        return (<Grid<Vocab>
            addRowHandler={() => setEditHandler(emptyEditData)}
            setEditHandler={setEditHandler}
            deleteHandler={deleteHandler}
            columns={[
                {title: '#', field: 'edit'},
                {title: 'German', field: 'german', width: '33%'},
                {title: 'Japanese', field: 'japanese', width: '33%'},
                {title: 'Kanji', field: 'kanji', width: '33%'}
            ]}
            data={vocabs}
        />);
    }, [vocabs]);
    const editModal = useMemo(() => {
        const cancelHandler = (): void => {
            setEditData(emptyEditData);
            setShowEditModal(false);
        };
        const onChangeHandler = (field: string, value: string): void => {
            if (editData && (field === 'german' || field === 'japanese' || field === 'kanji')) {
                const newEditData: Vocab = {...editData};
                newEditData[field] = value;
                setEditData(newEditData);
            }
        };
        const saveHandler = (): void => {
            if (editData) {
                post<Vocab, Vocab>('/vocab', editData, (data) => {
                    const foundedVocabs = vocabs.filter(vocab => vocab.id).filter(vocab => vocab.id !== data.id);
                    setVocabs([...foundedVocabs, data]);
                    setEditData(emptyEditData);
                });
            }
        };
        return (<VocabularyEditModal cancelHandler={cancelHandler}
                                     onChangeHandler={onChangeHandler}
                                     saveHandler={saveHandler}
                                     show={showEditModal}
                                     modalClosed={cancelHandler}
                                     editData={editData}
        />);
    }, [editData, showEditModal]);

    return (<>
        {editModal}
        {grid}
    </>);
};

export default VocabularyView;
