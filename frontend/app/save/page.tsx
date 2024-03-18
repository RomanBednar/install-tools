'use client'

import { useTransition } from 'react';
import { addComment } from '../actions/add-comment';

export default async function ArticleComment(props: any) {
    return (
        <form action={addComment}>
            <input type="text" name="username" value={props.username} />
            <input type="text" name="password" value={props.password}/>
            <button type="submit">Add Comment</button>
        </form>
    )
}