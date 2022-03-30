import React, { useEffect, useState } from 'react';
import './Root.css';
import Button from 'react-bootstrap/Button';
import userMap from '../applications/user.json';
import axios from 'axios';

type RootProps = Readonly<{}>;

const backendPort = 8080;
const host = `http:/127.0.0.1:${backendPort}`;

const Root: React.FC<RootProps> = () => {
    const [userChoice, setUserChoice] = useState('oda');
    const [user, setUser] = useState('');
    const [text_buffer, setBuffer] = useState('');
    const [result, setResult] = useState('');

    const handleUserChange = React.useCallback((event) => {
        setUserChoice(event.target.value);
        console.log(userChoice);
    }, [setUserChoice]);

    const handleUserSubmit = React.useCallback((event) => {
        event.preventDefault();
        if (userChoice === 'other') {
            // generate user
        } else if (userChoice === '') {
            console.log('Error')
        } else {
            setUser(userChoice);
        }
    }, [setUser, userChoice]);

    const handleCommandChange = React.useCallback((event) => {
        console.log(`Command: ${event.target.value}`);
        setBuffer(event.target.value)
    }, [setBuffer]);

    const handleCommandSubmit = React.useCallback((event) => {
        event.preventDefault();
        console.log(`com: ${text_buffer}`);
        // TODO: send request
        axios.get(`${host}/user/${user}/commands`)
            .then((res) => {
                if (res.status === 200) {
                    console.log(res);
                } else if (res.status === 204) {
                    console.log(res);
                }
            })
            .catch(console.error);
    }, [text_buffer, setResult])

    console.log(`userChoise is ${userChoice}`);
    console.log(`user is ${user}`);

    let sessionId = 0;
    let flag = false;
    for (const data of userMap.user) {
        if (data.name === user) {
            flag = true;
            // TODO: key-ex
            const pub = data['pub-key'];
            const sec = data['sec-key'];
            sessionId = 9;
            break;
        }
    }
    if (!flag) {
        console.log("Error: not implemented!")
    }

    return (
        <div className="Root">
            <header className="Root-header">
                <form>
                    <select value={userChoice} onChange={handleUserChange}>
                        <option value="oda">oda</option>
                        <option value="tada">tada</option>
                        <option value="other">新規</option>
                    </select>
                    <Button variant="primary" onClick={handleUserSubmit}>決定</Button>
                </form>
                <p>ユーザ名：{user}</p>
                <form>
                    <input type="text" onChange={handleCommandChange} />
                </form>
                <Button variant="primary" onClick={handleCommandSubmit}>実行</Button>
                <p>結果：{result}</p>
            </header>
        </div>
    )
}

export default Root;